module Manifest.Update exposing (..)

import Manifest.Models exposing (..)
import Manifest.Messages exposing (..)


update : Msg -> List Manifest -> ( List Manifest, Cmd Msg )
update message manifests =
    case message of
        FetchAllDone manifests' ->
            ( manifests', Cmd.none )

        FetchAllFail error ->
            ( manifests, Cmd.none )
