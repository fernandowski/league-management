import {View} from "react-native";
import {Surface, Text} from "react-native-paper";
import {useEffect} from "react";
import {SeasonStandings, SeasonStandingsResponse, useData} from "@/hooks/useData";
import TableList, {ColumnDefinition} from "@/components/TableList/TableList";


export interface SeasonStandingProps {
    seasonId: string
}

const columns: ColumnDefinition<SeasonStandings>[] = [

    {key: 'team_name', title: 'Name', width: 100},
    {key: 'total_points', title: 'Points', width: 80},
    {key: 'games_played', title: 'Games Played', width: 100},
    {key: 'total_wins', title: 'W'},
    {key: 'total_losses', title: 'L'},
    {key: 'total_ties', title: 'T'},
    {key: 'total_goals', title: 'G'},
];
export default function SeasonStanding(props: SeasonStandingProps) {

    const {fetchData, data, error, fetching} = useData<SeasonStandingsResponse>();


    useEffect(() => {
        if (props.seasonId) {
            fetchData(`/v1/seasons/${props.seasonId}/standings`)
        }
    }, [props.seasonId]);

    if (fetching) {
        return <View><Text>Fetching Details</Text></View>
    }

    return (
        <View>
            <Surface>
                {error && <Text>{error}</Text>}
                {data && <TableList data={data.standings} columns={columns}/>}
            </Surface>
        </View>
    )
}
