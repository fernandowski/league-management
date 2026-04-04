import {StyleSheet, View} from "react-native";
import {Modal, Text} from "react-native-paper";
import MembershipManagement from "@/components/League/MembershipManagement";

export interface LeagueMembershipModalProps {
    organizationId: string,
    leagueId: string | null,
    onDismiss: () => void,
    open: boolean
}

export default function LeagueMembershipModal(props: LeagueMembershipModalProps) {
    return (
        <Modal visible={props.open} contentContainerStyle={[styles.modal]} onDismiss={props.onDismiss}>
            <MembershipManagement leagueId={props.leagueId}/>
        </Modal>
    )
}


const styles = StyleSheet.create({
    modal: {
        flex: 1,
        padding: 16,
        backgroundColor: "white",
        width: "100%",
        maxHeight: '90%',
        maxWidth: '90%',
        zIndex: 200,
        alignSelf: "center"
    },
    container: {
        flex: 0.7,
        padding: 16
    }
})
