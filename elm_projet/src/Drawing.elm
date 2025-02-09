module Drawing exposing (..)

import Parsing exposing (Instruction(..))
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Basics exposing (degrees, cos, sin)

type alias Positions_turtle = 
    { x : Float
    , y : Float
    , angle : Float
    , couleur : String 
    }

-- Crée le vecteur svg à partir de deux points de coordonnes a et b
avancer : List Instruction -> Positions_turtle -> Positions_turtle -> (List (Svg msg)) -> (Positions_turtle, (List (Svg msg)))
avancer etapes position_a position_b dessin_svg =
    dessiner etapes position_b ( List.append dessin_svg
        [ line [ x1 (String.fromFloat position_a.x)
                , y1 (String.fromFloat position_a.y)
                , x2 (String.fromFloat position_b.x)
                , y2 (String.fromFloat position_b.y)
                , stroke position_a.couleur
                , strokeWidth "2"  -- Épaisseur du trait
                ] []
        ]

    )

-- pour répéter l'instruction n fois
repeat : Int -> List Instruction -> (Positions_turtle, List (Svg msg)) -> (Positions_turtle, List (Svg msg))
repeat n liste_instructions (positions_turtle, dessin_svg) =
    if n <= 0 then
        (positions_turtle, dessin_svg)
    else
        let
            (new_positions, new_svg) = dessiner liste_instructions positions_turtle dessin_svg
        in
        repeat (n - 1) liste_instructions (new_positions, new_svg)


-- donne positions et liste svg correspondant au dessin attendu avec la liste d'instrctions
dessiner : List Instruction -> Positions_turtle -> (List (Svg msg)) -> (Positions_turtle, (List (Svg msg)))
dessiner etapes positions_turtle dessin_svg =
    case etapes of
        -- y'a plus d'étapes à effectuer donc on renvoie le résultat
        [] -> 
            (positions_turtle, dessin_svg)

        -- y'a encore des instructions à dessiner
        (instruction :: liste_instructions) -> 
            case instruction of
                -- on avance
                Forward distance ->
                    avancer liste_instructions positions_turtle
                        (Positions_turtle (positions_turtle.x + cos (degrees positions_turtle.angle) * distance)
                                  (positions_turtle.y + sin (degrees positions_turtle.angle) * distance)
                                  positions_turtle.angle
                                  positions_turtle.couleur
                        )
                        dessin_svg
                -- on veut tourner à droite
                Right angle ->
                    dessiner liste_instructions (Positions_turtle positions_turtle.x positions_turtle.y (positions_turtle.angle + angle) positions_turtle.couleur) dessin_svg
                -- on veut tourner à gauche
                Left angle ->
                    dessiner liste_instructions (Positions_turtle positions_turtle.x positions_turtle.y (positions_turtle.angle - angle) positions_turtle.couleur) dessin_svg
                -- on veut répéter n fois
                Repeat nb instructions ->
                    repeat nb instructions (positions_turtle, dessin_svg)