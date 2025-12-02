{ pkgs, ... }: {
  channel = "stable-24.05";

  packages = [
    pkgs.go
    pkgs.air
    pkgs.postgresql
  ];

  env = {
    POSTGRESQL_CONN_STRING = "postgres://user:mypassword@localhost:5432/youtube?sslmode=disable";
  };

  idx = {
    workspace = {
      onCreate = {
        setup = ''
          # tạo data dir nếu chưa có
          mkdir -p pgdata

          # init db nếu chưa init
          if [ ! -f pgdata/PG_VERSION ]; then
            initdb -D pgdata -U user
          fi

          # tạo socket dir (Postgres dùng /tmp)
          mkdir -p /tmp

          # start postgres lần đầu để tạo database
          pg_ctl -D pgdata -o "-k /tmp" -l logfile start

          # chờ server start
          sleep 2

          # tạo database youtube
          psql -h /tmp -U user -d postgres -c "CREATE DATABASE youtube;"

          # tạo password cho user
          psql -h /tmp -U user -d postgres -c "ALTER USER user WITH PASSWORD 'mypassword';"

          # stop server tạm thời
          pg_ctl -D pgdata -l logfile stop
        '';
      };

      onStart = {
        start_pg = ''
          # đảm bảo socket /tmp tồn tại
          mkdir -p /tmp
          # start postgres mỗi lần workspace mở
          pg_ctl -D pgdata -o "-k /tmp" -l logfile start
        '';
      };
    };

    previews = {
      enable = true;
      previews = {
        web = {
          command = [
            "sh"
            "-c"
            "PORT=$PORT air -c ./air.toml"
          ];
          manager = "web";
        };
      };
    };

    extensions = [
      "golang.go"
    ];
  };
}
