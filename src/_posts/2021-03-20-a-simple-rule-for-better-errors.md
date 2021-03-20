---
title: >-
    A Simple Rule for Better Errors
description: >-
    ...and some examples of the rule in action.
tags: tech
---

This post will describe a simple rule for writing error messages that I've
been using for some time and have found to be worthwhile. Using this rule I can
be sure that my errors are propagated upwards with everything needed to debug
problems, while not containing tons of extraneous or duplicate information.

This rule is not specific to any particular language, pattern of error
propagation (e.g. exceptions, signals, simple strings), or method of embedding
information in errors (e.g. key/value pairs, formatted strings).

I do not claim to have invented this system, I'm just describing it.

## The Rule

Without more ado, here's the rule:

> A function sending back an error should not include information the caller
> could already know.

Pretty simple, really, but the best rules are. Keeping to this rule will result
in error messages which, once propagated up to their final destination (usually
some kind of logger), will contain only the information relevant to the error
itself, with minimal duplication.

The reason this rule works in tandem with good encapsulation of function
behavior. The caller of a function knows only the inputs to the function and, in
general terms, what the function is going to do with those inputs. If the
returned error only includes information outside of those two things then the
caller knows everything it needs to know about the error, and can continue on to
propagate that error up the stack (with more information tacked on if necessary)
or handle it in some other way.

## Examples

(For examples I'll use Go, but as previously mentioned this rule will be useful
in any other language as well.)

Let's go through a few examples, to show the various ways that this rule can
manifest in actual code.

**Example 1: Nothing to add**

In this example we have a function which merely wraps a call to `io.Copy` for
two files:

```go
func copyFile(dst, src *os.File) error {
	_, err := io.Copy(dst, src)
	return err
}
```

In this example there's no need to modify the error from `io.Copy` before
returning it to the caller. What would we even add? The caller already knows
which files were involved in the error, and that the error was encountered
during some kind of copy operation (since that's what the function says it
does), so there's nothing more to say about it.

**Example 2: Annotating which step an error occurs at**

In this example we will open a file, read its contents, and return them as a
string:

```go
func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("reading contents: %w", err)
	}

	return string(contents), nil
}
```

In this example there are two different steps which could result in an error:
opening the file and reading its contents. If an error is returned then our
imaginary caller doesn't know which step the error occurred at. Using our rule
we can infer that it would be good to annotate at _which_ step the error is
from, so the caller is able to have a fuller picture of what went wrong.

Note that each annotation does _not_ include the file path which was passed into
the function. The caller already knows this path, so an error being returned
back which reiterates the path is unnecessary.

**Example 3: Annotating which argument was involved**

In this example we will read two files using our function from example 2, and
return the concatenation of their contents as a string.

```go
func concatFiles(pathA, pathB string) (string, error) {
	contentsA, err := readFile(pathA)
	if err != nil {
		return "", fmt.Errorf("reading contents of %q: %w", pathA, err)
	}

	contentsB, err := readFile(pathB)
	if err != nil {
		return "", fmt.Errorf("reading contents of %q: %w", pathB, err)
	}

	return contentsA + contentsB, nil
}
```

Like in example 2 we annotate each error, but instead of annotating the action
we annotate which file path was involved in each error. This is because if we
simply annotated with the string `reading contents` like before it wouldn't be
clear to the caller _which_ file's contents couldn't be read. Therefore we
include which path the error is relevant to.

**Example 4: Layering**

In this example we will show how using this rule habitually results in easy to
read errors which contain all relevant information surrounding the error. Our
example reads one file, the "full" file, using our `readFile` function from
example 2. It then reads the concatenation of two files, the "split" files,
using our `concatFiles` function from example 3. It finally determines if the
two strings are equal:

```go
func verifySplits(fullFilePath, splitFilePathA, splitFilePathB string) error {
	fullContents, err := readFile(fullFilePath)
	if err != nil {
		return fmt.Errorf("reading contents of full file: %w", err)
	}

	splitContents, err := concatFiles(splitFilePathA, splitFilePathB)
	if err != nil {
		return fmt.Errorf("reading concatenation of split files: %w", err)
	}

	if fullContents != splitContents {
		return errors.New("full file's contents do not match the split files' contents")
	}

	return nil
}
```

As previously, we don't annotate the file paths for the different possible
errors, but instead say _which_ files were involved. The caller already knows
the paths, there's no need to reiterate them if there's another way of referring
to them.

Let's see what our errors actually look like! We run our new function using the
following:

```go
	err := verifySplits("full.txt", "splitA.txt", "splitB.txt")
	fmt.Println(err)
```

Let's say `full.txt` doesn't exist, we'll get the following error:

```
reading contents of full file: opening file: open full.txt: no such file or directory
```

The error is simple, and gives you everything you need to understand what went
wrong: while attempting to read the full file, during the opening of that file,
our code found that there was no such file. In fact, the error returned by
`os.Open` contains the name of the file, which goes against our rule, but it's
the standard library so what can ya do?

Now, let's say that `splitA.txt` doesn't exist, then we'll get this error:

```
reading concatenation of split files: reading contents of "splitA.txt": opening file: open splitA.txt: no such file or directory
```

Now we did include the file path here, and so the standard library's failure to
follow our rule is causing us some repitition. But overall, within the parts of
the error we have control over, the error is concise and gives you everything
you need to know what happened.

## Exceptions

As with all rules, there are certainly exceptions. The primary one I've found is
that certain helper functions can benefit from bending this rule a bit. For
example, if there is a helper function which is called to verify some kind of
user input in many places, it can be helpful to include that input value within
the error returned from the helper function:

```go
func verifyInput(str string) error {
    if err := check(str); err != nil {
        return fmt.Errorf("input %q was bad: %w", str, err)
    }
    return nil
}
```

`str` is known to the caller so, according to our rule, we don't need to include
it in the error. But if you're going to end up wrapping the error returned from
`verifyInput` with `str` at every call site anyway it can be convenient to save
some energy and break the rule. It's a trade-off, convenience in exchange for
consistency.

Another exception might be made with regards to stack traces.

In the set of examples given above I tended to annotate each error being
returned with a description of where in the function the error was being
returned from. If your language automatically includes some kind of stack trace
with every error, and if you find that you are generally able to reconcile that
stack trace with actual code, then it may be that annotating each error site is
unnecessary, except when annotating actual runtime values (e.g. an input
string).

As in all things with programming, there are no hard rules; everything is up to
interpretation and the specific use-case being worked on. That said, I hope what
I've laid out here will prove generally useful to you, in whatever way you might
try to use it.

