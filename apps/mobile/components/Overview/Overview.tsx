import AddOrganization from "@/components/Overview/AddOrganization";
import {OrganizationTable} from "@/components/Overview/OrganizationTable";
import ViewContent from "@/components/Layout/ViewContent";
import {StyleSheet, View} from "react-native";
import { AppText } from "@/components/ui/AppText";

export default function Overview() {
    return (
        <ViewContent>
            <View style={styles.header}>
                <AppText variant={"headlineMedium"}>Organizations</AppText>
                <AppText variant={"bodyLarge"}>
                    Start here by creating an organization, then select it before managing leagues, teams, and seasons.
                </AppText>
            </View>
            <AddOrganization/>
            <OrganizationTable/>
        </ViewContent>
    )
}

const styles = StyleSheet.create({
    header: {
        gap: 8,
        marginBottom: 24,
    },
});
