import {useOrganizationStore} from "@/stores/organizationStore";
import {useCallback, useState} from "react";
import {LeagueList} from "@/components/League/LeagueList";
import {Button, Portal} from "react-native-paper";
import {useFocusEffect, useRouter} from "expo-router";
import ViewContent from "@/components/Layout/ViewContent";
import {View} from "react-native";
import AddLeagueModal from "@/components/League/AddLeagueModal";

export default function LeagueOverview() {
    const {organization} = useOrganizationStore();
    const [openAddLeagueModal, setOpenAddLeagueModal] = useState<boolean>(false);
    const [showLeagueMembershipModal, setShowLeagueMembershipModal] = useState(false);
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
                setShowLeagueMembershipModal(false);
            };
        }, [organization])
    );

    return (
        <View>
            <ViewContent>
                <Button
                    style={[{alignSelf: 'flex-end'}]} mode={'elevated'}
                    onPress={handleOpenAddLeagueModal}>
                    + Add League
                </Button>
                <LeagueList refresh={refreshList}></LeagueList>
            </ViewContent>
            <Portal>
                <AddLeagueModal open={openAddLeagueModal} onSave={handleSaveLeague} onClose={onClose}/>
            </Portal>
        </View>
    )
}
