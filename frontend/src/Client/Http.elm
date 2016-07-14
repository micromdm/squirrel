module Client.Http exposing (..)

import Http exposing (Error(..))
import Task
import Json.Decode
import Manifest.Models exposing (Manifest, manifest)
import Manifest.Messages exposing (..)


fetchAllManifests : Cmd Msg
fetchAllManifests =
    Http.get (Json.Decode.list manifest) "/api/v1/manifests"
        |> Task.perform FetchAllFail FetchAllDone


reportError : Http.Error -> Http.Error
reportError error =
    case error of
        Http.Timeout ->
            Debug.log "API timeout" error

        Http.NetworkError ->
            Debug.log "Network error contacting API" error

        Http.UnexpectedPayload payload ->
            Debug.log ("Unexpected payload from API: " ++ payload) error

        Http.BadResponse status payload ->
            Debug.log ("Unexpected status/payload from API: " ++ (toString status) ++ "/" ++ payload) error
