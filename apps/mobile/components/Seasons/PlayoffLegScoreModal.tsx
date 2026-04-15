import StyledModal from "@/components/StyledModal";
import {Button, Text} from "react-native-paper";
import {DimensionValue, StyleSheet, useWindowDimensions, View} from "react-native";
import {useEffect, useState} from "react";

import {AppTextField} from "@/components/ui/AppTextField";
import {PlayoffTieMatchResponse} from "@/hooks/useData";
import {useAppTheme} from "@/theme/theme";

interface PlayoffLegScoreModalProps {
    isOpen: boolean
    match: PlayoffTieMatchResponse
    onClose: () => void
    onSave: (scores: {home: number; away: number}) => void
}

export default function PlayoffLegScoreModal(props: PlayoffLegScoreModalProps) {
    const [homeScore, setHomeScore] = useState(String(props.match.home_score));
    const [awayScore, setAwayScore] = useState(String(props.match.away_score));
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;
    const theme = useAppTheme();

    useEffect(() => {
        if (props.isOpen) {
            setHomeScore(String(props.match.home_score));
            setAwayScore(String(props.match.away_score));
        }
    }, [props.isOpen, props.match]);

    const modalWidth: DimensionValue = isLargeScreen ? "60%" : "90%";

    const handleScoreChange = (value: string, setter: React.Dispatch<React.SetStateAction<string>>) => {
        if (/^\d*$/.test(value)) {
            setter(value.replace(/^0+(?!$)/, ""));
        }
    };

    return (
        <StyledModal isOpen={props.isOpen} width={modalWidth}>
            <View style={{alignItems: "center"}}>
                <View style={styles.header}>
                    <Text style={[styles.title, {color: theme.colors.onSurface}]}>Match {props.match.match_order}</Text>
                    <Text style={[styles.subtitle, {color: theme.colors.onSurfaceVariant}]}>
                        {props.match.home_team} vs. {props.match.away_team}
                    </Text>
                </View>
                <View style={styles.body}>
                    <View style={[styles.scoreContainer, !isLargeScreen ? {flexDirection: "column"} : undefined]}>
                        <View style={styles.team}>
                            <Text style={[styles.teamName, {color: theme.colors.onSurface}]}>{props.match.home_team}</Text>
                            <AppTextField
                                style={[styles.scoreInput, {backgroundColor: theme.colors.surfaceVariant}]}
                                outlineStyle={{borderColor: theme.colors.outline, borderRadius: 10}}
                                value={homeScore}
                                keyboardType="numeric"
                                onChangeText={(value) => handleScoreChange(value, setHomeScore)}
                            />
                        </View>
                        <Text style={[styles.separator, {color: theme.colors.onSurfaceVariant}]}> - </Text>
                        <View style={styles.team}>
                            <Text style={[styles.teamName, {color: theme.colors.onSurface}]}>{props.match.away_team}</Text>
                            <AppTextField
                                style={[styles.scoreInput, {backgroundColor: theme.colors.surfaceVariant}]}
                                outlineStyle={{borderColor: theme.colors.outline, borderRadius: 10}}
                                value={awayScore}
                                keyboardType="numeric"
                                onChangeText={(value) => handleScoreChange(value, setAwayScore)}
                            />
                        </View>
                    </View>
                </View>
                <View style={styles.footer}>
                    <Button onPress={props.onClose}>Cancel</Button>
                    <Button onPress={() => props.onSave({home: Number(homeScore), away: Number(awayScore)})}>Save</Button>
                </View>
            </View>
        </StyledModal>
    );
}

const styles = StyleSheet.create({
    header: {
        marginBottom: 16,
        gap: 4,
        alignItems: "center",
    },
    title: {
        fontSize: 18,
        fontWeight: "bold",
        textAlign: "center",
    },
    subtitle: {
        fontSize: 14,
        textAlign: "center",
    },
    body: {
        marginBottom: 16,
    },
    scoreContainer: {
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
    team: {
        alignItems: "center",
    },
    teamName: {
        marginBottom: 8,
        width: 200,
        textAlign: "center",
    },
    scoreInput: {
        width: 60,
        height: 40,
        textAlign: "center",
    },
    separator: {
        fontSize: 24,
        fontWeight: "bold",
    },
    footer: {
        flexDirection: "row",
        justifyContent: "space-between",
    },
});
