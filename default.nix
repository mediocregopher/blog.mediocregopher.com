let
    utils = (import ./nix) {};
    pkgs = utils.pkgs;
    system = utils.system;
in
    rec {
        srv = (import ./srv).build;
        static = (import ./static).build;
        config = (import ./config.nix);

        service = pkgs.writeText "mediocregopher-mediocre-blog" ''
            [Unit]
            Description=mediocregopher mediocre blog
            Requires=network.target
            After=network.target

            [Service]
            Restart=always
            RestartSec=1s
            User=mediocregopher
            ExecStart=${srv}/bin/mediocre-blog \
                -pow-secret   "${config.powSecret}" \
                -ml-smtp-addr "${config.mlSMTPAddr}" \
                -ml-smtp-auth "${config.mlSMTPAuth}" \
                -data-dir     "${config.dataDir}" \
                -public-url   "${config.publicURL}" \
                -static-dir   "${static}" \
                -listen-proto "${config.listenProto}" \
                -listen-addr  "${config.listenAddr}"

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
