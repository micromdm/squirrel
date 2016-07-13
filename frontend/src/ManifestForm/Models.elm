module ManifestForm.Models exposing (..)

import Manifest.Models exposing (Manifest)
import Material


type alias Model =
    { mdl : Material.Model
    , manifestForm : Maybe Manifest
    , catalogs : List String
    }


initialModel : Model
initialModel =
    { mdl = Material.model
    , manifestForm = Nothing
    , catalogs = [ "production", "testing" ]
    }
