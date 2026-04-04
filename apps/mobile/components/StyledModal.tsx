import {Portal, Modal, Button} from "react-native-paper";
import {View, StyleSheet, DimensionValue} from "react-native";

interface ModalProps {
    isOpen: boolean;
    children: React.ReactNode;
    width?: DimensionValue;
    height?: DimensionValue
}

export default function StyledModal(props: ModalProps) {
    return (
        <Portal>
            <Modal
                visible={props.isOpen}
                dismissable={false}
                contentContainerStyle={[
                    styles.modal,
                    props.width ? {width: props.width} : styles.defaultWidth,
                    props.height ? {height: props.height} : styles.defaultHeight,
                ]}>
                <View style={styles.modalMainContent}>
                    {props.children}
                </View>
            </Modal>
        </Portal>
    );
}

const styles = StyleSheet.create({
    modal: {
        backgroundColor: "white",
        alignSelf: "center",
        borderRadius: 8,
        padding: 16,
        shadowColor: "#000",
        shadowOffset: {width: 0, height: 2},
        shadowOpacity: 0.25,
        shadowRadius: 4,
        elevation: 5,
    },
    modalMainContent: {
        flex: 1,
    },
    defaultWidth: {
        width: "auto",
    },
    defaultHeight: {
        height: "auto",
    }
});
