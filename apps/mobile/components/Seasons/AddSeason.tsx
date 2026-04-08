import {Surface} from "react-native-paper";
import {StyleSheet, View} from "react-native";
import AddSeasonModal from "@/components/Seasons/AddSeasonModal";
import React, {useState} from "react";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";


export interface AddSeasonProps {
    leagueId: string
    onSeasonAdded: () =>  void
}
export default function AddSeason(props: AddSeasonProps) {
    const [openModal, setOpenModal] = useState<boolean>(false);

    const handleSave = () => {
        setOpenModal(false);
        props.onSeasonAdded()
    };

    const openSeasonModal = () => {
        setOpenModal(!openModal);
    }

    return (
        <Surface style={styles.surfaceContainer} elevation={4}>
            <View style={[styles.buttonContainer]}>
                <AppText>No Season Configured</AppText>
                <AppButton mode={'elevated'} onPress={openSeasonModal}>+ Add Season</AppButton>
            </View>
            <AddSeasonModal onSave={handleSave} onClose={openSeasonModal} open={openModal} leagueId={props.leagueId}/>
        </Surface>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
    buttonContainer: {
        marginVertical: 27,
        alignItems: "center",
        gap: 8
    },
    surfaceContainer: {
        marginTop: 27,
        gap: 18,
        height: "80%"
    }

})
