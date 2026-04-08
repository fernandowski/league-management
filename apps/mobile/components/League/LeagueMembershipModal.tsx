import {StyleSheet} from "react-native";
import MembershipManagement from "@/components/League/MembershipManagement";
import StyledModal from "@/components/StyledModal";

export interface LeagueMembershipModalProps {
    organizationId: string,
    leagueId: string | null,
    onDismiss: () => void,
    open: boolean
}

export default function LeagueMembershipModal(props: LeagueMembershipModalProps) {
    return (
        <StyledModal
            isOpen={props.open}
            onDismiss={props.onDismiss}
            dismissable
            contentContainerStyle={styles.modal}
        >
            <MembershipManagement leagueId={props.leagueId}/>
        </StyledModal>
    )
}


const styles = StyleSheet.create({
    modal: {
        flex: 1,
        width: "100%",
        maxHeight: '90%',
        maxWidth: '90%',
        alignSelf: "center"
    },
    container: {
        flex: 0.7,
        padding: 16
    }
})
