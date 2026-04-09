import {StyleSheet, View} from "react-native";
import LeagueDropdown from "@/components/League/LeagueDropdown";
import {useLocalSearchParams, useRouter} from "expo-router";
import ViewContent from "@/components/Layout/ViewContent";
import {useEffect, useState} from "react";
import LeagueDetails from "@/components/League/LeagueDetails";
import { AppText } from "@/components/ui/AppText";

export default function Index() {
    const router = useRouter();

    const {id} = useLocalSearchParams();
    const [leagueId, setLeagueId] = useState<string>(Array.isArray(id) ? id[0] : id || "");
    const [refreshLeagueDetails] = useState(false);


    const onLeagueChange = (newLeagueId: string) => {
        if (newLeagueId !== leagueId) {
            router.setParams({
                id: newLeagueId
            });
            setLeagueId(newLeagueId);
        }
    }
    useEffect(() => {
        const selectedId = Array.isArray(id) ? id[0] : id || "";
        if (selectedId !== leagueId) {
            setLeagueId(selectedId);
        }
    }, [id, leagueId]);

    return (
        <ViewContent>
            <View style={{ flex: 1}}>
                <View style={styles.pageHeader}>
                    <AppText variant={"headlineMedium"}>League</AppText>
                    <AppText variant={"bodyLarge"}>
                        Select a league to manage details, membership, and season activity.
                    </AppText>
                </View>
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
    pageHeader: {
        gap: 8,
        marginBottom: 16,
    },
    dropdownContainer: {
        flexDirection: "row",
        paddingHorizontal: 16,
        marginBottom: 16,
    },
    detailsContainer: {
        flex: 1,
    },
});
