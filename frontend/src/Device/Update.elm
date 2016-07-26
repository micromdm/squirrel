module Device.Update exposing (..)

import Device.Models exposing (..)
import Device.Messages exposing (..)


update : Msg -> List Device -> ( List Device, Cmd Msg )
update message devices =
    case message of
        FetchDevicesDone devices' ->
            ( devices', Cmd.none )

        FetchDevicesFail error ->
            ( devices, Cmd.none )
