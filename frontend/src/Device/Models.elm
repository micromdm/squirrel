module Device.Models exposing (device, Device)

import Json.Decode exposing ((:=))
import Json.Decode.Extra exposing ((|:), withDefault)


type alias Device =
    { serialNumber : String
    , catalogs : List String
    , notes : String
    , user : String
    , hostName : String
    , templateManifest : String
    , depStatus : String
    }


device : Json.Decode.Decoder Device
device =
    Json.Decode.succeed Device
        |: (withDefault "" ("serial_nubmer" := Json.Decode.string))
        |: (withDefault [] ("catalogs" := Json.Decode.list Json.Decode.string))
        |: (withDefault "" ("notes" := Json.Decode.string))
        |: (withDefault "" ("user" := Json.Decode.string))
        |: (withDefault "" ("hostname" := Json.Decode.string))
        |: (withDefault "" ("template_manifest" := Json.Decode.string))
        |: (withDefault "" ("dep_status" := Json.Decode.string))
