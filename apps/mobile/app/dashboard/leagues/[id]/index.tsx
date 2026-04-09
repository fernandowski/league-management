import {StyleSheet, View} from "react-native";
import LeagueDropdown from "@/components/League/LeagueDropdown";
import {useLocalSearchParams, useRouter} from "expo-router";
import ViewContent from "@/components/Layout/ViewContent";
import {useEffect, useState} from "react";
import LeagueDetails from "@/components/League/LeagueDetails";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

export default function Index() {
    const router = useRouter();
    const theme = useAppTheme();

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
            <View style={styles.container}>
                <View style={styles.pageHeader}>
                    <AppText variant={"headlineMedium"}>League workspace</AppText>
                    <AppText variant={"bodyLarge"} style={{color: theme.colors.onSurfaceVariant}}>
                        Select a league to manage details, membership, and season activity.
                    </AppText>
                </View>
                <View style={styles.detailsContainer}>
                    <LeagueDetails
                        leagueId={leagueId}
                        refresh={refreshLeagueDetails}
                        selector={(
                            <LeagueDropdown onChange={onLeagueChange} selected={leagueId}/>
                        )}
                    />
                </View>
            </View>
        </ViewContent>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        gap: 16,
    },
    pageHeader: {
        gap: 8,
    },
    detailsContainer: {
        flex: 1,
    },
});
