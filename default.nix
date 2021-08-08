{

    pkgs ? import (fetchTarball {
        name = "nixpkgs-21-05";
        url = "https://github.com/NixOS/nixpkgs/archive/7e9b0dff974c89e070da1ad85713ff3c20b0ca97.tar.gz";
        sha256 = "1ckzhh24mgz6jd1xhfgx0i9mijk6xjqxwsshnvq789xsavrmsc36";
    }) {},

    baseConfig ? import ./config.nix,

}: rec {

    config = baseConfig // {
        redisListenPath = "${config.runDir}/redis";
    };

    static = (import ./static) { inherit pkgs; };

    srv = (import ./srv) {
        inherit pkgs config;
        staticBuild=static.build;
    };

    redisCfg = pkgs.writeText "mediocre-blog-redisCfg" ''
        port 0
        unixsocket ${config.redisListenPath}
        daemonize no
        loglevel notice
        logfile ""
        appendonly yes
        appendfilename "appendonly.aof"
        dir ${config.dataDir}/redis
    '';

    redisBin = pkgs.writeScript "mediocre-blog-redisBin" ''
        #!/bin/sh
        mkdir -p ${config.dataDir}/redis
        exec ${pkgs.redis}/bin/redis-server ${redisCfg}
    '';

    circusCfg = pkgs.writeText "mediocre-blog-circusCfg" ''
        [circus]
        endpoint = tcp://127.0.0.1:0
        pubsub_endpoint = tcp://127.0.0.1:0

        [watcher:srv]
        cmd = ${srv.bin}
        numprocesses = 1

        [watcher:redis]
        cmd = ${redisBin}
        numprocesses = 1
    '';

    entrypoint = pkgs.writeScript "mediocre-blog-entrypoint" ''
        #!/bin/sh
        mkdir -p ${config.runDir}
        mkdir -p ${config.dataDir}
        exec ${pkgs.circus}/bin/circusd ${circusCfg}
    '';

    service = pkgs.writeText "mediocre-blog" ''
        [Unit]
        Description=mediocregopher mediocre blog
        Requires=network.target
        After=network.target

        [Service]
        Restart=always
        RestartSec=1s
        User=mediocregopher
        ExecStart=${entrypoint}

        [Install]
        WantedBy=multi-user.target
    '';

    install = pkgs.writeScript "mediocre-blog" ''
        set -e -x

        sudo cp ${service} /etc/systemd/system/mediocregopher-mediocre-blog.service
        sudo systemctl daemon-reload
        sudo systemctl enable mediocregopher-mediocre-blog.service
        sudo systemctl restart mediocregopher-mediocre-blog.service
    '';
}
