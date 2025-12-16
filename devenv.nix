{ pkgs, lib, config, inputs, ... }:

{
  env = {
    TERN_MIGRATIONS = "./migrations";
  };

  packages = [
    pkgs.git
    pkgs.golangci-lint
    pkgs.tailwindcss_4
    pkgs.rustywind
    pkgs.just
  ];

  languages.go.enable = true;
}
