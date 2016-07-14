module Manifest.Messages exposing (..)

import Http
import Manifest.Models exposing (..)


type Msg
    = FetchAllDone (List Manifest)
    | FetchAllFail Http.Error
