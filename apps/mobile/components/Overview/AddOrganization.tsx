import {View, StyleSheet} from "react-native";
import {useState} from "react";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import Joi from "joi";
import {SubmitHandler, useForm} from "react-hook-form";
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

interface CreateOrganizationData {
    name: string
}


export default function AddOrganization() {
    const [showModal, setShowModal] = useState<boolean>(false);
    const {fetchOrganizations, setOrganization, organizations} = useOrganizationStore();
    const {control, handleSubmit, reset, formState: {errors, isSubmitting}} = useForm<CreateOrganizationData>(
        {
            resolver: joiResolver(schema),
            defaultValues: {
                name: "",
            }
        }
    );

    const handleShowModal = () => {
        setShowModal((currentState) => !currentState);
    }

    const handleOnSaveOrganization: SubmitHandler<CreateOrganizationData> = async (data: CreateOrganizationData): Promise<void> => {
        const response = await apiRequest('/v1/organizations', {
            method: 'POST',
            body: data
        });

        await fetchOrganizations();

        if (response?.id) {
            setOrganization(response.id);
        }

        reset();
        setShowModal(false);
    }

    return (
        <View style={[styles.container, organizations.length > 0 && styles.compactContainer]}>
            <AppText style={styles.description}>
                {organizations.length > 0
                    ? "Create another organization to manage a different league group."
                    : "Create your first organization to start managing leagues, teams, and seasons."}
            </AppText>
            <AppButton mode={"contained"} onPress={handleShowModal}>
                Add Organization
            </AppButton>
            <StyledModal isOpen={showModal} onDismiss={handleShowModal} width={"90%"} height={"auto"}>
                <View style={[styles.formContainer]}>
                    <AppText variant={"titleMedium"}>Create Organization</AppText>
                    <ControlledTextInput label='Name' name={'name'} control={control} error={errors.name?.message}/>
                </View>
                <View style={styles.formActionButtons}>
                    <AppButton onPress={handleShowModal}>Cancel</AppButton>
                    <AppButton loading={isSubmitting} onPress={handleSubmit(handleOnSaveOrganization)}>Save</AppButton>
                </View>
            </StyledModal>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        alignItems: "flex-start",
        gap: 8,
        marginBottom: 24,
    },
    compactContainer: {
        justifyContent: "flex-start",
    },
    description: {
        maxWidth: 540,
    },
    formContainer: {
        gap: 16,
    },
    formActionButtons: {
        alignSelf: "flex-end",
        flexDirection: "row",
    },
})
