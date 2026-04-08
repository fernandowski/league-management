import {useCallback, useState} from "react";
import {LeagueList} from "@/components/League/LeagueList";
import {useFocusEffect} from "expo-router";
import ViewContent from "@/components/Layout/ViewContent";
import {View} from "react-native";
import AddLeagueModal from "@/components/League/AddLeagueModal";
import { AppButton } from "@/components/ui/AppButton";

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
                <AppButton
                    style={[{alignSelf: 'flex-end'}]} mode={'elevated'}
                    onPress={handleOpenAddLeagueModal}>
                    + Add League
                </AppButton>
                <LeagueList refresh={refreshList}></LeagueList>
            </ViewContent>
            <AddLeagueModal open={openAddLeagueModal} onSave={handleSaveLeague} onClose={onClose}/>
        </View>
    )
}
