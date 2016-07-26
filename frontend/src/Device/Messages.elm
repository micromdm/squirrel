module Device.Messages exposing (..)

import Http
import Device.Models exposing (..)


type Msg
    = FetchDevicesDone (List Device)
    | FetchDevicesFail Http.Error
