module ManifestForm.Messages exposing (..)

import Http
import Manifest.Models exposing (..)
import Material


type Action
    = Add
    | Remove


type Msg
    = FetchManifest String Manifest
    | FetchManifestFail Http.Error
    | Mdl Material.Msg
    | SetCatalog String Bool
    | SetDisplayName String
    | SetUser String
    | SetNotes String
    | SetManagedInstalls String Action
    | SetOptionalInstalls String Action
    | SetManagedUninstalls String Action
    | SetIncludedManifests String Action
    | ClearForm
    | NoOp
