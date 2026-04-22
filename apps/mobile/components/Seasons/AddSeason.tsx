import {StyleSheet, View} from "react-native";
import AddSeasonModal from "@/components/Seasons/AddSeasonModal";
import React, {useState} from "react";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";
import { AppCard } from "@/components/ui/AppCard";
import { useAppTheme } from "@/theme/theme";


export interface AddSeasonProps {
    leagueId: string
    onSeasonAdded: () =>  void
}
export default function AddSeason(props: AddSeasonProps) {
    const [openModal, setOpenModal] = useState<boolean>(false);
    const theme = useAppTheme();

    const handleSave = () => {
        setOpenModal(false);
        props.onSeasonAdded()
    };

    const openSeasonModal = () => {
        setOpenModal(!openModal);
    }

    return (
        <AppCard style={[styles.surfaceContainer, {borderColor: theme.colors.outline}]}>
            <AppCard.Content style={styles.content}>
                <View style={styles.copy}>
                    <AppText variant="headlineSmall">No active season</AppText>
                    <AppText variant="bodyLarge" style={{color: theme.colors.onSurfaceVariant}}>
                        Create a season to unlock scheduling, standings, and match management for this league.
                    </AppText>
                </View>
                <View style={[styles.actionRow, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                    <View style={styles.actionCopy}>
                        <AppText variant="titleMedium">Create a season</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            Start with a season name, then plan the schedule when your teams are ready.
                        </AppText>
                    </View>
                    <AppButton variant="submit" onPress={openSeasonModal}>Add Season</AppButton>
                </View>
            </AppCard.Content>
            <AddSeasonModal onSave={handleSave} onClose={openSeasonModal} open={openModal} leagueId={props.leagueId}/>
        </AppCard>
    )
}

const styles = StyleSheet.create({
    surfaceContainer: {
        borderRadius: 24,
    },
    content: {
        gap: 16,
        paddingVertical: 8,
    },
    copy: {
        gap: 6,
    },
    actionRow: {
        borderWidth: 1,
        borderRadius: 20,
        padding: 16,
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        gap: 12,
        flexWrap: "wrap",
    },
    actionCopy: {
        flex: 1,
        minWidth: 220,
        gap: 4,
    },
})
