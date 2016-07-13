module Main exposing (..)

import Html exposing (..)
import Navigation
import String
import UrlParser exposing (Parser, (</>), format, int, oneOf, string)


-- Update


type Msg
    = NoOp


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case Debug.log "msg" msg of
        NoOp ->
            ( model, Cmd.none )



-- View


view : Model -> Html Msg
view model =
    case model.page of
        Home ->
            viewHome model


viewHome : Model -> Html Msg
viewHome model =
    div [] [ text "hello world" ]



-- Model


type alias Model =
    { page : Page
    }



-- Routing


type Page
    = Home


toHash : Page -> String
toHash page =
    case page of
        Home ->
            "#home"


urlUpdate : Result String Page -> Model -> ( Model, Cmd Msg )
urlUpdate result model =
    case Debug.log "result" result of
        Err _ ->
            ( model, Navigation.modifyUrl (toHash model.page) )

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
    urlUpdate result initialModel


initialModel : Model
initialModel =
    { page = Home
    }
