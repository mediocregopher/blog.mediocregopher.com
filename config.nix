{
  runDir = "/tmp/mediocre-blog/run";
  dataDir = "/tmp/mediocre-blog/data";

  powSecret = "ssshhh";
  mlSMTPAddr = "";
  mlSMTPAuth = "";
  mlPublicURL = "http://localhost:4000";
  listenProto = "tcp";
  listenAddr = ":4000";

  # If empty then a derived static directory is used
  staticProxyURL = "http://127.0.0.1:4002";

  # password is "bar". This should definitely be changed for prod.
  httpAuthUsers = {
    "foo" = "$2a$13$0JdWlUfHc.3XimEMpEu1cuu6RodhUvzD9l7iiAqa4YkM3mcFV5Pxi";
  };
}
