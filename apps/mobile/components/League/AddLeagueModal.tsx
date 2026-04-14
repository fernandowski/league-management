import {StyleSheet, View} from "react-native";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import Joi from "joi";
import {useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";
import {apiRequest} from "@/api/api";
import {useOrganizationStore} from "@/stores/organizationStore";
import StyledModal from "@/components/StyledModal";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";

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
        } catch {

        }
    }

    return (
        <StyledModal isOpen={props.open} onDismiss={props.onClose}>
            <View style={[styles.formContainer, styles.formFields]}>
                <View>
                    <AppText variant="titleMedium">League Name</AppText>
                    <ControlledTextInput label='Name' name={'name'} control={control} error={errors.name?.message}/>
                </View>
            </View>
            <View style={styles.formActionButtons}>
                <AppButton variant="secondary" onPress={props.onClose}>Close</AppButton>
                <AppButton variant="submit" onPress={handleSubmit(handleSave)}>Save</AppButton>
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
