import { useState } from 'react';
import { StyleProp, StyleSheet, View, ViewStyle } from 'react-native';
import Joi from 'joi';
import { SubmitHandler, useForm } from 'react-hook-form';
import { joiResolver } from '@hookform/resolvers/joi';

import { apiRequest } from '@/api/api';
import ControlledTextInput from '@/components/FormControls/ControlledTextInput';
import StyledModal from '@/components/StyledModal';
import { AppButton } from '@/components/ui/AppButton';
import { AppText } from '@/components/ui/AppText';
import { useOrganizationStore } from '@/stores/organizationStore';


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

interface AddOrganizationProps {
    buttonLabel?: string;
    description?: string;
    modalTitle?: string;
    style?: StyleProp<ViewStyle>;
}

export default function AddOrganization({
    buttonLabel = "Add Organization",
    description,
    modalTitle = "Create Organization",
    style,
}: AddOrganizationProps) {
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
        <View style={[styles.container, organizations.length > 0 && styles.compactContainer, style]}>
            {(description !== "" || description === undefined) && (
                <AppText style={styles.description}>
                    {description ?? (organizations.length > 0
                        ? "Create another organization to manage a different league group."
                        : "Create your first organization to start managing leagues, teams, and seasons.")}
                </AppText>
            )}
            <AppButton variant="submit" onPress={handleShowModal}>
                {buttonLabel}
            </AppButton>
            <StyledModal isOpen={showModal} onDismiss={handleShowModal} width={"90%"} height={"auto"}>
                <View style={[styles.formContainer]}>
                    <AppText variant={"titleMedium"}>{modalTitle}</AppText>
                    <ControlledTextInput label='Name' name={'name'} control={control} error={errors.name?.message}/>
                </View>
                <View style={styles.formActionButtons}>
                    <AppButton variant="secondary" onPress={handleShowModal}>Cancel</AppButton>
                    <AppButton variant="submit" loading={isSubmitting} onPress={handleSubmit(handleOnSaveOrganization)}>Save</AppButton>
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
