import {StyleSheet, View} from "react-native";
import {useOrganizationStore} from "@/stores/organizationStore";
import {useAppTheme} from "@/theme/theme";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";

export function OrganizationTable() {
    const {organization: selectedOrganizationId, organizations, setOrganization} = useOrganizationStore();
    const theme = useAppTheme();

    if (organizations.length === 0) {
        return null;
    }

    return (
        <View style={[styles.outerContainer]}>
            <View style={[styles.viewContainer]}>
                {organizations.map((organization) =>
                    <AppCard
                        style={[
                            styles.card,
                            organization.id === selectedOrganizationId && styles.selectedCard,
                            organization.id === selectedOrganizationId && {
                                borderColor: theme.colors.primary,
                                backgroundColor: theme.colors.primaryContainer,
                            },
                        ]}
                        key={organization.id}
                        onPress={() => setOrganization(organization.id)}
                    >
                        <AppCard.Content>
                            <AppText variant={"titleMedium"}>{organization.name}</AppText>
                            <AppText variant={"bodySmall"} style={{color: theme.colors.onSurfaceVariant}}>
                                {organization.id === selectedOrganizationId
                                    ? "Selected"
                                    : "Tap to manage this organization"}
                            </AppText>
                        </AppCard.Content>
                    </AppCard>)}
            </View>
        </View>
    )
}


const styles = StyleSheet.create({
    outerContainer: {
        gap: 12,
    },
    viewContainer: {
        gap: 8,
    },
    card: {
        marginTop: 8,
        borderWidth: 1,
    },
    selectedCard: {
        borderWidth: 2,
    },
})
