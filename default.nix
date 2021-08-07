let
    utils = (import ./nix) {};
    pkgs = utils.pkgs;
    system = utils.system;
in
    {config ? ./config.nix}: rec {
        config = (import ./config.nix);

        static = (import ./static).build;

        srv = (import ./srv).build;
        srvBin = pkgs.writeScript "mediocregopher-mediocre-blog-srvBin" ''
            #!/bin/sh
            exec ${srv}/bin/mediocre-blog \
                -pow-secret   "${config.powSecret}" \
                -ml-smtp-addr "${config.mlSMTPAddr}" \
                -ml-smtp-auth "${config.mlSMTPAuth}" \
                -data-dir     "${config.dataDir}" \
                -public-url   "${config.publicURL}" \
                -static-dir   "${static}" \
                -listen-proto "${config.listenProto}" \
                -listen-addr  "${config.listenAddr}"
        '';

        redisCfg = pkgs.writeText "mediocregopher-mediocre-blog-redisCfg" ''
            port 0
            unixsocket ${config.redisListenPath}
            daemonize no
            loglevel notice
            logfile ""
            appendonly yes
            appendfilename "appendonly.aof"
            dir ${config.dataDir}/redis
        '';

        redisBin = pkgs.writeScript "mediocregopher-mediocre-blog-redisBin" ''
            #!/bin/sh
            mkdir -p ${config.dataDir}/redis
            exec ${pkgs.redis}/bin/redis-server ${redisCfg}
        '';

        circusCfg = pkgs.writeText "mediocregopher-mediocre-blog-circusCfg" ''
            [circus]
            endpoint = tcp://127.0.0.1:0
            pubsub_endpoint = tcp://127.0.0.1:0

            [watcher:srv]
            cmd = ${srvBin}
            numprocesses = 1

            [watcher:redis]
            cmd = ${redisBin}
            numprocesses = 1
        '';

        circusBin = pkgs.writeScript "mediocregopher-mediocre-blog-circusBin" ''
            exec ${pkgs.circus}/bin/circusd ${circusCfg}
        '';

        service = pkgs.writeText "mediocregopher-mediocre-blog" ''
            [Unit]
            Description=mediocregopher mediocre blog
            Requires=network.target
            After=network.target

            [Service]
            Restart=always
            RestartSec=1s
            User=mediocregopher
            ExecStart=${circusBin}

            [Install]
            WantedBy=multi-user.target
        '';

        install = pkgs.writeScript "mediocregopher-mediocre-blog" ''
                set -e -x

                sudo cp ${service} /etc/systemd/system/mediocregopher-mediocre-blog.service
                sudo systemctl daemon-reload
                sudo systemctl enable mediocregopher-mediocre-blog.service
                sudo systemctl restart mediocregopher-mediocre-blog.service
        '';
    }
