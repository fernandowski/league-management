import {StyleSheet, View} from "react-native";
import {Button, Modal, Text} from "react-native-paper";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import {useState} from "react";
import Joi from "joi";
import {useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";
import {apiRequest} from "@/api/api";
import {useOrganizationStore} from "@/stores/organizationStore";
import StyledModal from "@/components/StyledModal";

const schema = Joi.object({
    name: Joi.string()
        .required()
        .messages({
            "string.empty": "Organization name is required",
        })
});

interface CreateLeagueData {
    name: string;
}

export interface AddLeagueModalProps {
    onSave: () => void;
    onClose: () => void;
    open: boolean;
}

export default function AddLeagueModal(props: AddLeagueModalProps) {
    const {organization} = useOrganizationStore();
    const {control, handleSubmit, formState: {errors}} = useForm<CreateLeagueData>(
        {
            resolver: joiResolver(schema),
        }
    );

    const handleSave = async (data: CreateLeagueData) => {
        try {
            await apiRequest('/v1/leagues', {
                method: 'POST',
                body: {
                    name: data.name,
                    organization_id: organization
                }
            })
            props.onSave();
        } catch (e) {

        }
    }

    return (
        <StyledModal isOpen={props.open}>
            <View style={[styles.formContainer, styles.formFields]}>
                <View>
                    <Text>Team Name</Text>
                    <ControlledTextInput label='Name' name={'name'} control={control} error={errors.name?.message}/>
                </View>
            </View>
            <View style={styles.formActionButtons}>
                <Button onPress={props.onClose}>Close</Button>
                <Button onPress={handleSubmit(handleSave)}>Save</Button>
            </View>
        </StyledModal>
    )
}

const styles = StyleSheet.create({
    formContainer: {
        justifyContent: "space-between",
        flex: 1
    },
    formFields: {
        maxWidth: 500,
    },
    formActionButtons: {
        alignSelf: "flex-end",
        flexDirection: "row"
    }
})


