module Client.Http exposing (..)

import Http exposing (Error(..))
import Task
import Json.Decode
import Manifest.Models exposing (Manifest, manifest)
import Manifest.Messages exposing (..)
import Device.Models exposing (Device, device)
import Device.Messages exposing (..)
import ManifestForm.Messages exposing (..)


fetchAllManifests : Cmd Manifest.Messages.Msg
fetchAllManifests =
    Http.get (Json.Decode.list manifest) "/api/v1/manifests"
        |> Task.perform FetchManifestsFail FetchManifestsDone


fetchManifest : String -> Cmd ManifestForm.Messages.Msg
fetchManifest name =
    Http.get manifest ("/api/v1/manifests/" ++ name)
        |> Task.perform FetchManifestFail (FetchManifest name)


fetchAllDevices : Cmd Device.Messages.Msg
fetchAllDevices =
    Http.get (Json.Decode.list device) "/api/v1/devices"
        |> Task.perform FetchDevicesFail FetchDevicesDone


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
