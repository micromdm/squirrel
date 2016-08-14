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

        SetManagedInstalls install action ->
            let
                removeItem manifest =
                    -- unused for now
                    { manifest | managedInstalls = List.filter (\i -> i /= install) manifest.managedInstalls }

                form' =
                    -- unused for now
                    Maybe.map (removeItem) model.manifestForm

                uncheckedAdd a =
                    a :: model.uManagedInstalls

                uncheckedRemove a =
                    List.filter (\i -> i /= install) model.uManagedInstalls

                uManagedInstalls' =
                    case action of
                        Add ->
                            uncheckedRemove install

                        Remove ->
                            uncheckedAdd install
            in
                { model
                    | manifestForm = model.manifestForm
                    , uManagedInstalls = uManagedInstalls'
                }
                    ! []

        SetOptionalInstalls install action ->
            let
                removeItem manifest =
                    { manifest | optionalInstalls = List.filter (\i -> i /= install) manifest.optionalInstalls }

                form' =
                    Maybe.map (removeItem) model.manifestForm
            in
                { model | manifestForm = form' } ! []

        SetManagedUninstalls install action ->
            let
                removeItem manifest =
                    { manifest | managedUninstalls = List.filter (\i -> i /= install) manifest.managedUninstalls }

                form' =
                    Maybe.map (removeItem) model.manifestForm
            in
                { model | manifestForm = form' } ! []

        SetIncludedManifests included action ->
            let
                removeItem manifest =
                    { manifest | includedManifests = List.filter (\i -> i /= included) manifest.includedManifests }

                form' =
                    Maybe.map (removeItem) model.manifestForm
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
