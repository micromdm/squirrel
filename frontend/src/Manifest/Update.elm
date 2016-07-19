module Manifest.Update exposing (..)

import Manifest.Models exposing (..)
import Manifest.Messages exposing (..)


update : Msg -> List Manifest -> ( List Manifest, Cmd Msg )
update message manifests =
    case message of
        FetchManifestsDone manifests' ->
            ( manifests', Cmd.none )

        FetchManifestsFail error ->
            ( manifests, Cmd.none )
