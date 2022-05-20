{
  runDir = "/tmp/mediocre-blog/run";
  dataDir = "/tmp/mediocre-blog/data";
  publicURL = "http://localhost:4000";

  powSecret = "ssshhh";
  mlSMTPAddr = "";
  mlSMTPAuth = "";
  httpListenProto = "tcp";
  httpListenAddr = ":4000";

  # password is "bar". This should definitely be changed for prod.
  httpAuthUsers = {
    "foo" = "$2a$13$0JdWlUfHc.3XimEMpEu1cuu6RodhUvzD9l7iiAqa4YkM3mcFV5Pxi";
  };

  # Very low, should be increased for prod.
  httpAuthRatelimit = "1s";
}
