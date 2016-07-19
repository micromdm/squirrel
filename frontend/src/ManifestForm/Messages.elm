module ManifestForm.Messages exposing (..)

import Http
import Manifest.Models exposing (..)
import Material


type Msg
    = FetchManifest String Manifest
    | FetchManifestFail Http.Error
    | Mdl Material.Msg
    | SetCatalog String Bool
    | SetDisplayName String
    | SetUser String
    | SetNotes String
    | ClearForm
    | NoOp
