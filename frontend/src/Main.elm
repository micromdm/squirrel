module Main exposing (..)

import Html exposing (..)
import Html.App as App
import Navigation
import String
import UrlParser exposing (Parser, (</>), format, int, oneOf, string)
import Manifest.Models exposing (Manifest)
import ManifestForm.Models as ManifestForm
import Manifest.Messages
import Manifest.Update
import ManifestForm.Messages
import ManifestForm.Update
import ManifestForm.View as ManifestForm
import Client.Http exposing (..)
import Material
import Material.Layout as Layout
import Material.Options as Options exposing (css, when)
import Material.Table as Table
import Material.Icon as Icon
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)


-- Update


type Msg
    = NoOp
    | Mdl Material.Msg
    | SelectTab Int
    | ManifestMsg Manifest.Messages.Msg
    | ManifestFormMsg ManifestForm.Messages.Msg
    | EditManifest String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case Debug.log "msg" msg of
        NoOp ->
            ( model, Cmd.none )

        ManifestMsg subMsg ->
            let
                ( manifests', cmd ) =
                    Manifest.Update.update subMsg model.manifests
            in
                ( { model | manifests = manifests' }, Cmd.map ManifestMsg cmd )

        ManifestFormMsg subMsg ->
            let
                ( manifest', cmd ) =
                    ManifestForm.Update.update subMsg model.manifestForm
            in
                ( { model | manifestForm = manifest' }, Cmd.map ManifestFormMsg cmd )

        EditManifest id ->
            let
                navCmd =
                    Navigation.newUrl (toHash <| ManifestForm id)

                fetchCmd =
                    Cmd.map ManifestFormMsg <| fetchManifest id
            in
                { model
                    | page = (ManifestForm id)
                }
                    ! [ Cmd.batch [ navCmd, fetchCmd ] ]

        Mdl msg ->
            Material.update Mdl msg model

        SelectTab k ->
            let
                tabToPage =
                    case k of
                        0 ->
                            ManifestList

                        _ ->
                            Home
            in
                { model | page = tabToPage } ! [ Navigation.newUrl (toHash tabToPage) ]


viewPage : Model -> Html Msg
viewPage model =
    case model.page of
        Home ->
            div []
                [ text "hello world" ]

        ManifestList ->
            viewManifestTable model

        ManifestForm _ ->
            App.map ManifestFormMsg <| ManifestForm.view model.manifestForm


view : Model -> Html Msg
view model =
    Layout.render Mdl
        model.mdl
        [ Layout.fixedHeader
        , Layout.waterfall True
        , Layout.onSelectTab SelectTab
        ]
        { header = header model
        , drawer = []
        , tabs = ( tabTitles, [] )
        , main = [ viewPage model ]
        }


header : Model -> List (Html Msg)
header model =
    [ Layout.row
        [ Options.nop
        , css "transition" "height 333ms ease-in-out 0s"
        ]
        [ Layout.title [] [ text "Munki Admin" ]
        , Layout.spacer
        , Layout.navigation []
            [ Layout.link [ Layout.href "https://github.com/micromdm/squirrel" ]
                [ span [] [ text "github" ] ]
            , Layout.link [ Layout.href "http://package.elm-lang.org/packages/debois/elm-mdl/latest/" ]
                [ text "Help" ]
            ]
        ]
    ]


tabTitles : List (Html a)
tabTitles =
    [ text "Manifests"
    , text "Catalogs"
    , text "Devices"
    , text "Groups"
    , text "Users"
    ]


viewManifestTable : Model -> Html Msg
viewManifestTable model =
    let
        firstCatalog catalogs =
            catalogs
                |> List.head
                |> Maybe.withDefault ""

        manifestRow manifest =
            Table.tr []
                [ Table.td [] [ text manifest.name ]
                , Table.td [] [ text <| Maybe.withDefault "" manifest.displayName ]
                , Table.td [] [ text <| firstCatalog manifest.catalogs ]
                , Table.td [] [ span [ class "manifest-edit-button", onClick <| EditManifest manifest.name ] [ Icon.i "edit" ] ]
                ]

        manifestItems =
            List.map (manifestRow) model.manifests
    in
        Table.table []
            [ Table.thead []
                [ Table.tr []
                    [ Table.th []
                        [ text "Name" ]
                    , Table.th []
                        [ text "Display Name" ]
                    , Table.th []
                        [ text "Catalogs" ]
                    , Table.th []
                        [ text "Edit" ]
                    ]
                ]
            , Table.tbody [] manifestItems
            ]



-- Model


type alias Model =
    { mdl : Material.Model
    , page : Page
    , manifests : List Manifest
    , manifestForm : ManifestForm.Model
    }



-- Routing


type Page
    = Home
    | ManifestList
    | ManifestForm String


toHash : Page -> String
toHash page =
    case page of
        Home ->
            "#home"

        ManifestList ->
            "#manifests"

        ManifestForm id ->
            "#manifests/" ++ id


urlUpdate : Result String Page -> Model -> ( Model, Cmd Msg )
urlUpdate result model =
    case Debug.log "result" result of
        Err _ ->
            ( model, Navigation.modifyUrl (toHash model.page) )

        Ok ((ManifestForm id) as page) ->
            { model
                | page = page
            }
                ! [ Cmd.map ManifestFormMsg <| fetchManifest id ]

        Ok page ->
            { model
                | page = page
            }
                ! []


hashParser : Navigation.Location -> Result String Page
hashParser location =
    UrlParser.parse identity pageParser (String.dropLeft 1 location.hash)


pageParser : Parser (Page -> a) a
pageParser =
    oneOf
        [ format Home (UrlParser.s "home")
        , format ManifestForm (UrlParser.s "manifests" </> string)
        , format ManifestList (UrlParser.s "manifests")
        ]



-- Main


main : Program Never
main =
    Navigation.program (Navigation.makeParser hashParser)
        { init = init
        , update = update
        , view = view
        , subscriptions = (\_ -> Sub.none)
        , urlUpdate = urlUpdate
        }


init : Result String Page -> ( Model, Cmd Msg )
init result =
    let
        ( model, routeMsg ) =
            urlUpdate result initialModel

        manifestsMsg =
            Cmd.map ManifestMsg fetchAllManifests
    in
        ( model, Cmd.batch [ routeMsg, manifestsMsg ] )


initialModel : Model
initialModel =
    { mdl = Material.model
    , page = Home
    , manifests = []
    , manifestForm = ManifestForm.initialModel
    }
