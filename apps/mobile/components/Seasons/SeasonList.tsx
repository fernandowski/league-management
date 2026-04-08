import {View} from "react-native";
import SeasonSearchBar from "@/components/Seasons/SeasonSearchBar";
import {useState} from "react";
import SeasonsTable from "@/components/Seasons/SeasonsTable";

export default function SeasonList() {
    const [leagueId, setLeagueId] = useState<null | string>(null)
    const onLeagueSelection = (value: string) => {
        setLeagueId(value);
    };

    return (
        <View>
            <View>
                <SeasonSearchBar onLeagueSelection={onLeagueSelection}/>
            </View>
            <View>
                <SeasonsTable leagueId={leagueId || ''}></SeasonsTable>
            </View>
        </View>
    )
}
