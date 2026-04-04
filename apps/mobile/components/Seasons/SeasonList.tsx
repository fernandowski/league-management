import {View, StyleSheet} from "react-native";
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

const styles = StyleSheet.create({
    container: {
        flex: 1,
        width: '80%',
        marginTop: 16,
        marginLeft: 16,
        gap: 16,
        justifyContent: "center"
    },
    searchContainer: {
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-between",
    },
    outerContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center'
    },
    viewContainer: {
        flex: 1,
        width: '80%',
        marginTop: 16,
        gap: 16
    },
})
