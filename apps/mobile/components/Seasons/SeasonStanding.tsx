import {View} from "react-native";
import {Surface} from "react-native-paper";
import {useEffect} from "react";
import {SeasonStandings, SeasonStandingsResponse, useData} from "@/hooks/useData";
import TableList, {ColumnDefinition} from "@/components/TableList/TableList";
import { AppText } from "@/components/ui/AppText";


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
    }, [fetchData, props.seasonId]);

    if (fetching) {
        return <View><AppText>Fetching Details</AppText></View>
    }

    return (
        <View>
            <Surface>
                {error && <AppText>{error}</AppText>}
                {data && <TableList data={data.standings} columns={columns}/>}
            </Surface>
        </View>
    )
}
