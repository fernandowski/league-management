import ViewContent from "@/components/Layout/ViewContent";
import SeasonList from "@/components/Seasons/SeasonList";
import { AppText } from "@/components/ui/AppText";
import { StyleSheet, View } from "react-native";

export default function Index() {
    return (
        <ViewContent>
            <View style={styles.pageHeader}>
                <AppText variant={"headlineMedium"}>Seasons</AppText>
                <AppText variant={"bodyLarge"}>
                    Find the season for a league and manage its schedule and progress.
                </AppText>
            </View>
            <SeasonList/>
        </ViewContent>
    )
}

const styles = StyleSheet.create({
    pageHeader: {
        gap: 8,
        marginBottom: 20,
    },
});
