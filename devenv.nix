{ pkgs, config, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";
  env = {
    DB_NAME = "test";
    DB_USER = "postgres";
    DB_PASS = "postgres";
    DB_DSN = "${config.env.DB_USER}:${config.env.DB_PASS}@localhost/${config.env.DB_NAME}?sslmode=disable";
  };

  # https://devenv.sh/packages/
  packages = [ pkgs.git ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/languages/
  languages.go.enable = true;
  languages.go.package = pkgs.go_1_21;

  services.postgres = {
    enable = true;
    package = pkgs.postgresql_15;
    initialDatabases = [{ name = config.env.DB_NAME; }];
    listen_addresses = "localhost";
    port = 5432;
    initialScript = ''
        CREATE USER ${config.env.DB_USER} SUPERUSER;
        ALTER USER ${config.env.DB_USER} WITH PASSWORD '${config.env.DB_PASS}';
    '';
  };

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
