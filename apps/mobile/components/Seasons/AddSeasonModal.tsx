import {StyleSheet, View} from "react-native";
import {Button, Modal, Text} from "react-native-paper";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import {useOrganizationStore} from "@/stores/organizationStore";
import {useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";
import Joi from "joi";
import {apiRequest} from "@/api/api";
import {useEffect} from "react";

export interface AddSeasonModalProps {
    onSave: () => void
    open: boolean
    leagueId: string
    onClose: () => void
}

const schema = Joi.object({
    name: Joi.string()
        .required()
        .messages({
            "string.empty": "Season Name name is required",
        })
});

interface CreateSeasonData {
    name: string
}
export default function AddSeasonModal(props: AddSeasonModalProps) {
    const {reset, control, handleSubmit, formState: {errors}} = useForm<CreateSeasonData>(
        {
            resolver: joiResolver(schema),
        }
    );

    const handleSave = async (data: CreateSeasonData) => {
        await apiRequest(`/v1/leagues/${props.leagueId}/seasons`, {method: "POST", body: {name: data.name}});
        reset();
        props.onSave();
    };

    const onClose = () => {
        reset();
        props.onClose();
    }

    return (
        <Modal visible={props.open} dismissable={false} contentContainerStyle={[styles.modal]}>
            <View style={[styles.formContainer]}>
                <View style={{gap: 16}}>
                    <Text>Season Name: </Text>
                    <ControlledTextInput  name={'name'} control={control} error={errors.name?.message}/>
                </View>
                <View style={{flexDirection: "row", justifyContent: "flex-end"}}>
                    <Button onPress={onClose}>Cancel</Button>
                    <Button onPress={handleSubmit(handleSave)}>Save</Button>
                </View>
            </View>
        </Modal>
    )
}

const styles = StyleSheet.create({
    modal: {
        flex: 0.8,
        padding: 16,
        backgroundColor: "white",
        width: "80%",
        maxHeight: 400,
        maxWidth: 400,
        zIndex: 200,
        alignSelf: "center"
    },
    formContainer: {
        flex: 1,
        padding: 16,
        justifyContent: "space-between"
    },
})
