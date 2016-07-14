module Manifest.Models exposing (manifest, newManifest, Manifest)

import Json.Decode exposing ((:=))
import Json.Decode.Extra exposing ((|:))


type alias Manifest =
    { name :
        String
    , catalogs :
        Maybe (List String)
    , displayName :
        Maybe String
    , includedManifests :
        Maybe (List String)
    , notes :
        Maybe String
    , user :
        Maybe String
    , managedInstalls :
        Maybe (List String)
    , managedUninstalls :
        Maybe (List String)
    , optionalInstalls :
        Maybe (List String)
    }


newManifest : Manifest
newManifest =
    { name = ""
    , catalogs = Nothing
    , displayName = Nothing
    , includedManifests = Nothing
    , notes = Nothing
    , user = Nothing
    , managedInstalls = Nothing
    , managedUninstalls = Nothing
    , optionalInstalls = Nothing
    }


manifest : Json.Decode.Decoder Manifest
manifest =
    Json.Decode.succeed Manifest
        |: ("filename" := Json.Decode.string)
        |: (Json.Decode.maybe ("catalogs" := Json.Decode.list Json.Decode.string))
        |: (Json.Decode.maybe ("display_name" := Json.Decode.string))
        |: (Json.Decode.maybe ("include_manifests" := Json.Decode.list Json.Decode.string))
        |: (Json.Decode.maybe ("notes" := Json.Decode.string))
        |: (Json.Decode.maybe ("user" := Json.Decode.string))
        |: (Json.Decode.maybe ("managed_installs" := Json.Decode.list Json.Decode.string))
        |: (Json.Decode.maybe ("managed_uninstalls" := Json.Decode.list Json.Decode.string))
        |: (Json.Decode.maybe ("optional_installs" := Json.Decode.list Json.Decode.string))
