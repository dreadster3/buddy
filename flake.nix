{
  description = "Buddy flake";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-23.11"; };

  outputs = { self, nixpkgs }:
    let
      # System types to support.
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      # Get the current git revision
      version = self.rev or self.dirtyRev;

    in {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          buddy = pkgs.buildGoModule {
            pname = "buddy";
            inherit version;

            ldflags = "-X main.Version=${version}";

            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;

            vendorHash = "sha256-9ZuzJEBi6DegS6kYzgiY7ZRnMUO7lZ+Ze/V7hZQJ3So=";
          };

          default = self.packages.${system}.buddy;
        });

      # Add dependencies that are only needed for development
      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });
    };
}
