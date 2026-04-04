import {StyleSheet, View} from "react-native";
import LeagueDropdown from "@/components/League/LeagueDropdown";
import {useLocalSearchParams, useRouter} from "expo-router";
import ViewContent from "@/components/Layout/ViewContent";
import {useEffect, useState} from "react";
import LeagueDetails from "@/components/League/LeagueDetails";

export default function Index() {
    const router = useRouter();

    const {id} = useLocalSearchParams();
    const [leagueId, setLeagueId] = useState<string>(Array.isArray(id) ? id[0] : id || "");
    const [refreshLeagueDetails, setRefreshLeagueDetails] = useState(false);


    const onLeagueChange = (newLeagueId: string) => {
        if (newLeagueId !== leagueId) {
            router.setParams({
                id: newLeagueId
            });
            setLeagueId(newLeagueId);
        }
    }

    const onMemberRefresh = () => {
        setRefreshLeagueDetails(!refreshLeagueDetails);
    }

    useEffect(() => {
        const selectedId = Array.isArray(id) ? id[0] : id || "";
        if (selectedId !== leagueId) {
            setLeagueId(selectedId);
        }
    }, [id]);

    return (
        <ViewContent>
            <View style={{ flex: 1}}>
                <View>
                    <View style={styles.dropdownContainer}>
                        <LeagueDropdown onChange={onLeagueChange} selected={leagueId}/>
                    </View>
                </View>
                <View style={styles.detailsContainer}>
                    <LeagueDetails  leagueId={leagueId} refresh={refreshLeagueDetails}/>
                </View>
            </View>
        </ViewContent>
    )
}

const styles = StyleSheet.create({
    dropdownContainer: {
        flexDirection: "row",
        paddingHorizontal: 16,
        marginBottom: 16,
    },
    detailsContainer: {
        flex: 1,
    },
});
