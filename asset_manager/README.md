# Asset manager
The asset manager package manages the assets of the games:

- This package deals with assets from all versions of the game.
- It can also deal with newly created assets that are not in the database yet.

To make it simple, all the assets have to follow the asset path rule:

  - All asset must have asset path.
  - The asset path is unique within version, so the same version can't have 2 assets with the same asset path.
  - The asset path is shared in all versions, so the same asset path should reference equivalent asset in different version. An asset path is allowed to exist in some version(s) but not the others.
  - The rule apply to newly created assets as well:

    - Newly created assets are assigned an unique asset path upon being inserted.

This package will load asset based on the asset path, no matter if it is a new asset, or an existing asset.

## New assets
Newly created assets are managed differently compared to existing assets:

- Instead of referencing by files and keys, new assets will reference by file.
- Instead of being stored based on language and platform, new assets will just be the asset path and the file names for various version.
- The actual files can be changed, but the change might not take effect until this package is reloaded.

Users can trigger a asset packing procedure which will:

- Pack the file on disk into an encrypted pack clients can use:
  
  - The files are packed based on the set of language / platform they are used in.
  - There is no way to manually group the files, for now.
- Generate the .sql files to insert the reference for the assets and packs into the various asset databases.
- Users should avoid working on unrelated assets to avoid packing unnecessary files. 