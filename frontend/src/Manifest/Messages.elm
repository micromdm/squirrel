module Manifest.Messages exposing (..)

import Http
import Manifest.Models exposing (..)


type Msg
    = FetchManifestsDone (List Manifest)
    | FetchManifestsFail Http.Error
