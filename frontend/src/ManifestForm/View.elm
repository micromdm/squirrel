module ManifestForm.View exposing (..)

import ManifestForm.Models exposing (..)
import ManifestForm.Messages exposing (..)
import Material.Progress as Loading
import Material.Textfield as Textfield
import Material.Toggles as Toggles
import Material.Button as Button
import Html exposing (..)
import Html.Attributes exposing (..)


-- import Html.Events exposing (onClick)


view : Model -> Html Msg
view model =
    case model.manifestForm of
        Nothing ->
            div [] [ Loading.indeterminate ]

        Just manifest ->
            let
                name =
                    Textfield.render Mdl
                        [ 0 ]
                        model.mdl
                        [ Textfield.label "Name"
                        , Textfield.floatingLabel
                        , Textfield.value manifest.name
                        , Textfield.disabled
                        ]

                displayName =
                    Textfield.render Mdl
                        [ 1 ]
                        model.mdl
                        [ Textfield.label "Display Name"
                        , Textfield.floatingLabel
                        , Textfield.onInput SetDisplayName
                        , Textfield.value <| Maybe.withDefault "" manifest.displayName
                        ]

                user =
                    Textfield.render Mdl
                        [ 2 ]
                        model.mdl
                        [ Textfield.label "User"
                        , Textfield.floatingLabel
                        , Textfield.onInput SetUser
                        , Textfield.value <| Maybe.withDefault "" manifest.user
                        ]

                notes =
                    Textfield.render Mdl
                        [ 3 ]
                        model.mdl
                        [ Textfield.label "Notes"
                        , Textfield.floatingLabel
                        , Textfield.value <| Maybe.withDefault "" manifest.notes
                        , Textfield.textarea
                        , Textfield.onInput SetNotes
                        , Textfield.rows 6
                        ]

                hasCatalog name =
                    List.member name manifest.catalogs

                catalog idx catalogName =
                    Toggles.switch Mdl
                        [ 4 + idx ]
                        model.mdl
                        [ Toggles.value (hasCatalog catalogName)
                        , Toggles.onClick <| SetCatalog catalogName (hasCatalog catalogName)
                        ]
                        [ text catalogName ]

                catalogs =
                    List.indexedMap (catalog) model.catalogs

                managedInstall idx install =
                    Toggles.checkbox Mdl
                        [ 500 + idx ]
                        model.mdl
                        [ Toggles.value True ]
                        [ text install ]

                managedInstalls =
                    List.indexedMap (managedInstall) manifest.managedInstalls

                save =
                    Button.render Mdl
                        [ 1000 ]
                        model.mdl
                        [ Button.raised
                        , Button.ripple
                        , Button.colored
                        , Button.onClick ClearForm
                        ]
                        [ text "Save" ]

                cancel =
                    Button.render Mdl
                        [ 1001 ]
                        model.mdl
                        [ Button.raised
                        , Button.ripple
                        , Button.colored
                        , Button.onClick ClearForm
                        ]
                        [ text "Cancel" ]
            in
                div [ class "manifestForm" ]
                    [ div [ class "manifestFormGeneral" ]
                        [ name
                        , displayName
                        , user
                        , notes
                        , text "Catalogs"
                        , div [] catalogs
                        , div [ class "manifestFormButtons" ] [ save, cancel ]
                        ]
                    , div [ class "manifestFormAdvanced" ]
                        [ text "Managed Installs"
                        , div [] managedInstalls
                        ]
                    ]
