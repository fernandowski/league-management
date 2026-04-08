import {StyleSheet, View} from "react-native";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import {useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";
import Joi from "joi";
import {apiRequest} from "@/api/api";
import StyledModal from "@/components/StyledModal";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";

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
        <StyledModal
            isOpen={props.open}
            onDismiss={onClose}
            contentContainerStyle={styles.modal}
        >
            <View style={[styles.formContainer]}>
                <View style={{gap: 16}}>
                    <AppText variant="titleMedium">Season Name</AppText>
                    <ControlledTextInput  name={'name'} control={control} error={errors.name?.message}/>
                </View>
                <View style={{flexDirection: "row", justifyContent: "flex-end"}}>
                    <AppButton onPress={onClose}>Cancel</AppButton>
                    <AppButton onPress={handleSubmit(handleSave)}>Save</AppButton>
                </View>
            </View>
        </StyledModal>
    )
}

const styles = StyleSheet.create({
    modal: {
        width: "80%",
        maxHeight: 400,
        maxWidth: 400,
    },
    formContainer: {
        flex: 1,
        padding: 16,
        justifyContent: "space-between"
    },
})
