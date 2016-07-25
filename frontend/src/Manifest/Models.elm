module Manifest.Models exposing (manifest, newManifest, Manifest)

import Json.Decode exposing ((:=))
import Json.Decode.Extra exposing ((|:), withDefault)


type alias Manifest =
    { name :
        String
    , catalogs :
        List String
    , displayName :
        Maybe String
    , includedManifests :
        List String
    , notes :
        Maybe String
    , user :
        Maybe String
    , managedInstalls :
        List String
    , managedUninstalls :
        List String
    , optionalInstalls :
        List String
    }


newManifest : Manifest
newManifest =
    { name = ""
    , catalogs = []
    , displayName = Nothing
    , includedManifests = []
    , notes = Nothing
    , user = Nothing
    , managedInstalls = []
    , managedUninstalls = []
    , optionalInstalls = []
    }


manifest : Json.Decode.Decoder Manifest
manifest =
    Json.Decode.succeed Manifest
        |: (withDefault "" ("filename" := Json.Decode.string))
        |: (withDefault [] ("catalogs" := Json.Decode.list Json.Decode.string))
        |: (Json.Decode.maybe ("display_name" := Json.Decode.string))
        |: (withDefault [] ("included_manifests" := Json.Decode.list Json.Decode.string))
        |: (Json.Decode.maybe ("notes" := Json.Decode.string))
        |: (Json.Decode.maybe ("user" := Json.Decode.string))
        |: (withDefault [] ("managed_installs" := Json.Decode.list Json.Decode.string))
        |: (withDefault [] ("managed_uninstalls" := Json.Decode.list Json.Decode.string))
        |: (withDefault [] ("optional_installs" := Json.Decode.list Json.Decode.string))
