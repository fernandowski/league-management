import {useCallback, useState} from "react";
import {LeagueList} from "@/components/League/LeagueList";
import {useFocusEffect} from "expo-router";
import ViewContent from "@/components/Layout/ViewContent";
import {StyleSheet, View} from "react-native";
import AddLeagueModal from "@/components/League/AddLeagueModal";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";

export default function LeagueOverview() {
    const [openAddLeagueModal, setOpenAddLeagueModal] = useState<boolean>(false);
    const [refreshList, setRefreshList] = useState<boolean>(false);

    const handleOpenAddLeagueModal = (): void => {
        setOpenAddLeagueModal(!openAddLeagueModal);
    }

    const handleSaveLeague = (): void => {
        setOpenAddLeagueModal(!openAddLeagueModal);
        setRefreshList(!refreshList);
    }

    const onClose = () => {
        setOpenAddLeagueModal(false);
    }


    useFocusEffect(
        useCallback(() => {
            return () => {
                setOpenAddLeagueModal(false);
            };
        }, [])
    );

    return (
        <View>
            <ViewContent>
                <View style={styles.pageHeader}>
                    <View style={styles.pageHeaderTopRow}>
                        <AppText variant="headlineMedium">Leagues</AppText>
                        <AppButton
                            style={styles.addButton}
                            mode={'contained'}
                            onPress={handleOpenAddLeagueModal}>
                            Add League
                        </AppButton>
                    </View>
                    <AppText variant="bodyLarge">
                        Create and manage leagues for the selected organization.
                    </AppText>
                </View>
                <LeagueList refresh={refreshList}></LeagueList>
            </ViewContent>
            <AddLeagueModal open={openAddLeagueModal} onSave={handleSaveLeague} onClose={onClose}/>
        </View>
    )
}

const styles = StyleSheet.create({
    pageHeader: {
        gap: 8,
        marginBottom: 20,
    },
    pageHeaderTopRow: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        gap: 12,
    },
    addButton: {
        marginBottom: 0,
    },
});
