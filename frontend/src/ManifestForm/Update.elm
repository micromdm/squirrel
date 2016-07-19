module ManifestForm.Update exposing (..)

import Manifest.Models exposing (..)
import ManifestForm.Models exposing (..)
import ManifestForm.Messages exposing (..)
import Navigation
import Material


update : Msg -> Model -> ( Model, Cmd Msg )
update message model =
    case Debug.log "manifest" message of
        FetchManifest name manifest' ->
            ( { model | manifestForm = Just (withFilename manifest' name) }, Cmd.none )

        FetchManifestFail error ->
            ( { model | manifestForm = Nothing }, Cmd.none )

        SetCatalog name disable ->
            let
                removeCatalog manifest =
                    { manifest | catalogs = List.filter (\e -> e /= name) manifest.catalogs }

                addCatalog manifest =
                    { manifest | catalogs = name :: manifest.catalogs }

                form' =
                    case disable of
                        True ->
                            Maybe.map (removeCatalog) model.manifestForm

                        False ->
                            Maybe.map (addCatalog) model.manifestForm
            in
                { model | manifestForm = form' } ! []

        SetDisplayName name ->
            let
                form' =
                    Maybe.map (\m -> { m | displayName = Just name }) model.manifestForm
            in
                { model | manifestForm = form' } ! []

        SetUser name ->
            let
                form' =
                    Maybe.map (\m -> { m | user = Just name }) model.manifestForm
            in
                { model | manifestForm = form' } ! []

        SetNotes notes ->
            let
                form' =
                    Maybe.map (\m -> { m | notes = Just notes }) model.manifestForm
            in
                { model | manifestForm = form' } ! []

        ClearForm ->
            ( { model | manifestForm = Nothing }, Navigation.newUrl "#manifests" )

        Mdl msg ->
            Material.update Mdl msg model

        NoOp ->
            model ! []


withFilename : Manifest -> String -> Manifest
withFilename manifest filename =
    { manifest | name = filename }
